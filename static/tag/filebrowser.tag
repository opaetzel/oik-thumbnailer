<filebrowser>
    <div class="filebrowser">
        <form onsubmit={ gotoDir }>
            <b>Current:</b><input ref="current" value={current} onkeyup={changeCurrent} class="statusbar">
        </form>
        <div class="files-table">
        <table class="pure-table pure-table-striped full-width-table">
            <tbody>
                <tr class="dir">
                    <td class="shrink" onclick={ parentDir } style="cursor: pointer;">
                        <i class="fa fa-arrow-up"></i>
                    </td>
                    <td class="expand">
                        <span onclick={ parentDir }>../</span>
                    </td>
                </tr>
                <tr each={ files } class=" {selected?'selected':''} {IsDir?'dir':''}">
                    <td class="shrink">
                        <i class="fa {IsDir?'fa-folder':'fa-file-o'}"></i>
                    </td>
                    <td class="expand">
                        <span onclick={ enterDir }> { Name }</span>
                    </td>
                </th>
            </tbody>
        </table>
        </div>
        <div if={!selectDir} class="filename">
            File name: <input value={filename} onkeyup={changeFilename}>
        </div>
    </div>
        

  <!-- this script tag is optional -->
  <script>
    this.files = opts.Files
    this.apiEndpoint = opts.apiEndpoint
    this.selectDir = opts.selectDir
    this.current = opts.Current
    this.filename = opts.defaultFilename
    this.changed = opts.changed

    if(this.filename) {
        this.changed(this.current + "/" + this.filename)
    } else {
        this.changed(this.current)
    } 


    changeFilename(e) {
        this.filename = e.target.value
    }

    changeCurrent(e) {
        this.current = e.target.value
    }

    selectFile(e) {
        var item = e.item
        if(this.selectDir && !item.IsDir) {
            return
        } else if (!this.selectDir && item.IsDir) {
            return
        }
        var oldVal = item.selected
        for(var otherItem of this.files) {
            otherItem.selected = false;
        }
        item.selected = !oldVal
    }

    gotoDir(e) {
        e.preventDefault()
        var file = {Name: this.current}
        axios.post(this.apiEndpoint + 'currentDir', file ).then((response) => {
            this.opts = response.data
            this.files = response.data.Files
            this.current = this.opts.Current
            if(this.filename) {
                this.changed(this.current + "/" + this.filename)
            } else {
                this.changed(this.current)
            } 
            this.update()
        })
        .catch(function(error) {
            console.log(error)
        })
    }

    enterDir(e) {
        if(!e.item.IsDir) {
            return;
        }
        axios.post(this.apiEndpoint + 'currentDir', e.item ).then((response) => {
            this.opts = response.data
            this.files = response.data.Files
            this.current = this.opts.Current
            if(this.filename) {
                this.changed(this.current + "/" + this.filename)
            } else {
                this.changed(this.current)
            }
            this.update()
        })
        .catch(function(error) {
            console.log(error)
        })
    }

    parentDir() {
        axios.get(this.apiEndpoint + 'parentDir').then((response) => {
            this.opts = response.data
            this.files = response.data.Files
            this.current = this.opts.Current
            if(this.filename) {
                this.changed(this.current + "/" + this.filename)
            } else {
                this.changed(this.current)
            } 
            this.update()
        })
        .catch(function(error) {
            console.log(error)
        })
    }

</script>

</filebrowser>
