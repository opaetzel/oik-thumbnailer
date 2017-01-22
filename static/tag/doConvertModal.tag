<doConvertModal>
    <div id="myModal" class="modal {visible?'visible':''}">

        <!-- Modal content -->
        <div class="modal-content">
            <div class="modal-header">
                <span class="close" onclick={dismiss}>&times;</span>
                <h2 if={ !done && !working }>Convert?</h2>
                <h2 if={working}>Converting...</h2>
                <h2 if={done}>Done!</h2>
            </div>
            <div class="modal-body">
                <div if={ !done && !working }>
                    <b>Input folder: </b>{conv.sourceDir}<br>
                    <b>Target file: </b>{conv.destFile}<br>
                    <b>Width: </b>{conv.width}<br>
                    <b>Height: </b>{conv.height}
                </div>
                <div if={working}>
                    <center>
                        <img src="css/2.gif">
                    </center>
                </div>
                <div if={done && !working }>
                    Files packed in {conv.destFile}:
                    <ul>
                        <li each={packedFiles}>
                            {Name}
                        </li>
                    </ul>
                </div>
                <div if={error}>
                    An error occured...
                </div>
            </div>
            <div class="modal-footer">
                <div class="pure-g">
                    <div class="pure-u-1">
                        <div class="controls pull-right">
                            <span if={!working && !done}>
                                <button class="pure-button" onclick={ dismiss }>Cancel</button>
                                <button class="pure-button" onclick={ doConvert }>Convert</button>
                            </span>
                            <span if={done}
                                <button class="pure-button" onclick={ dismiss }>OK</button>
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>


    <script>
        this.conv = opts
        this.visible = false
        this.done = false

        opts.modalControl.on('show', () => {
            this.visible = true
            this.update()
        })

        doConvert() {
            this.working = true
            this.update()
            axios.put('api/convert', this.conv).then((response) => {
                console.log(response.data);
                this.working = false
                this.done = true
                this.packedFiles = response.data
                this.update()
            }).catch((error) => {
                console.log(error)
                this.error = true
                this.working = false
                this.update()
            });
        }
        
        dismiss() {
            this.visible = false
            this.update
        }
</script>

</doConvertModal>
