package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

type ByEndNum []os.FileInfo

func (s ByEndNum) Len() int {
	return len(s)
}
func (s ByEndNum) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByEndNum) Less(i, j int) bool {
	splitI := strings.Split(s[i].Name(), ".")
	splitJ := strings.Split(s[j].Name(), ".")
	if len(splitI) < 2 || len(splitJ) < 2 {
		return s[i].Name() < s[j].Name()
	}
	splitI2 := strings.Split(splitI[0], "_")
	splitJ2 := strings.Split(splitJ[0], "_")
	maybeNumI := splitI2[len(splitI2)-1]
	maybeNumJ := splitJ2[len(splitJ2)-1]

	numI, err := strconv.Atoi(maybeNumI)
	if err != nil {
		return s[i].Name() < s[j].Name()
	}
	numJ, err := strconv.Atoi(maybeNumJ)
	if err != nil {
		return s[i].Name() < s[j].Name()
	}
	return numI < numJ
}

func createPackage(input, output, dimensions *string) {
	width, height, err := getDim(dimensions)
	if err != nil {
		flag.Usage()
		return
	}
	_, _ = createPackageWNum(input, output, width, height)
}

func createPackageWNum(input, output *string, width, height int) ([]MyFile, error) {
	fmt.Printf("width: %d height: %d\n", width, height)
	files, err := ioutil.ReadDir(*input)
	if err != nil {
		fmt.Printf("Could not read input directory: %s\n", *input)
	}
	sort.Sort(ByEndNum(files))
	out, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Could not create file: %s\n", *output)
		return nil, err
	}
	defer out.Close()
	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	packedFiles := make([]MyFile, 0)
	for _, f := range files {
		if f.Name() == ".DS_STORE" {
			continue
		}
		fmt.Println(f.Name())
		var img image.Image
		file, err := os.Open(*input + "/" + f.Name())
		if err != nil {
			fmt.Printf("Could not read file: %s\n", f.Name())
			continue
		}
		defer file.Close()
		if strings.HasSuffix(f.Name(), ".jpg") || strings.HasSuffix(f.Name(), ".jpeg") {
			img, err = jpeg.Decode(file)
		} else if strings.HasSuffix(f.Name(), ".png") {
			img, err = png.Decode(file)
		} else {
			fmt.Printf("Unknown file format: %s\n", f.Name())
			continue
		}
		thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.NearestNeighbor)

		err = addFile(tw, f, &thumbnail)
		if err != nil {
			return nil, err
		}
		packedFiles = append(packedFiles, MyFile{Name: f.Name(), Size: f.Size(), IsDir: false})
	}
	return packedFiles, nil
}

func addFile(tw *tar.Writer, fi os.FileInfo, thumbnail *image.Image) error {
	// now lets create the header as needed for this file within the tarball
	var buf bytes.Buffer
	jpeg.Encode(&buf, *thumbnail, nil)
	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	header.Size = int64(buf.Len())
	// write the header to the tarball archive
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	// copy the file data to the tarball
	_, err = io.Copy(tw, &buf)
	if err != nil {
		return err
	}
	return nil
}

func getDim(dimensions *string) (width, height int, err error) {
	if *dimensions == "" {
		return 0, 0, errors.New("no dimensions specified")
	}

	split := strings.Split(*dimensions, "x")
	if len(split) != 2 {
		return 0, 0, errors.New("no x in dimensions")
	}

	width, err = strconv.Atoi(split[0])
	if err != nil {
		return
	}

	height, err = strconv.Atoi(split[1])
	if err != nil {
		return
	}

	return
}
