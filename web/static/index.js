"use strict"

document.addEventListener("alpine:init", () => {
    Alpine.data("global", () => ({
        text: "fdsfdsfds",
    }))

    Alpine.data("uploadState", () => ({
        filesToUpload: [],
        filesUploadingCount: 0,

        formatBytes(bytes, decimals = 2) {
            if (bytes === 0) return "0 Bytes"
            const k = 1024
            const dm = decimals < 0 ? 0 : decimals
            const sizes = [
                "Bytes",
                "KB",
                "MB",
                "GB",
                "TB",
                "PB",
                "EB",
                "ZB",
                "YB",
            ]
            const i = Math.floor(Math.log(bytes) / Math.log(k))
            return (
                parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) +
                " " +
                sizes[i]
            )
        },

        handleAdd(files) {
            // check that files have unique name
            const expectedLength = this.filesToUpload.length + files.length
            const allFiles = this.filesToUpload.concat(files)
            const actualLength = new Set(allFiles.map((file) => file.name)).size
            if (expectedLength !== actualLength) {
                console.log("Files with the same name are not allowed!")
                return
            }
            this.filesToUpload = allFiles
            console.log("added - ", files.length)
        },

        handleDrop(event) {
            event.preventDefault()
            this.handleAdd(Array.from(event.dataTransfer.files))
        },

        handleSelect(event) {
            this.handleAdd(Array.from(event.target.files))
        },

        handleRemove(filename) {
            this.filesToUpload = this.filesToUpload.filter(
                (file) => file.name !== filename
            )
        },

        handleUploadFile(file) {
            console.log("handle upload file", file)

            let xhr = new XMLHttpRequest()

            xhr.onreadystatechange = () => {
                console.log(`onreadystatechange: ${xhr.readyState}`)
            }

            xhr.onload = () => {
                console.log(
                    "onload (got response from server): ",
                    xhr.status,
                    xhr.response,
                    xhr.statusText
                )
            }

            xhr.onerror = () => {
                console.log("onerror")
            }

            xhr.onprogress = (event) => {
                console.log(
                    `onprogress: loaded: ${event.loaded}, total: ${event.total}, server sent Content-Length? ${event.lengthComputable}`
                )
            }

            xhr.upload.onprogress = (event) => {
                console.log(
                    `upload.onprogress: loaded: ${event.loaded}, total: ${event.total}, so progress is ${event.loaded / event.total}`
                )
            }

            xhr.upload.onload = () => {
                console.log(`upload.onloadend: successfully uploaded`)
            }

            xhr.onloadend = () => {
                if (xhr.status === 202) {
                    console.log("onloadend: success")
                } else {
                    console.log("onloadend: fail " + this.status)
                }

                this.filesUploadingCount--
                if (this.filesUploadingCount === 0) {
                    this.filesToUpload = []
                }
            }

            // xhr.timeout = 10000 // timeout 10 seconds, wait for server response
            xhr.open("POST", `/api/upload/${file.name}`)
            xhr.send(file)
        },

        handleBatchUpload() {
            if (this.filesUploadingCount > 0) return
            this.filesUploadingCount = this.filesToUpload.length

            this.filesToUpload.forEach((file) => this.handleUploadFile(file))
        },

        // async handleUploadDeprecated() {
        //     if (this.isUploading) return
        //
        //     console.log("handleUpload():", this.filesToUpload[0])
        //     this.isUploading = true
        //
        //     const formData = new FormData()
        //     this.filesToUpload.forEach((file) => formData.append("files", file))
        //     const path = "/api/upload"
        //
        //     try {
        //         const response = await fetch(path, {
        //             method: "POST",
        //             body: formData,
        //             headers: {
        //                 "X-Requested-With": "XMLHttpRequest", // To identify AJAX request on the server
        //             },
        //         })
        //
        //         if (!response.ok) {
        //             const errorData = await response.text()
        //             console.log(
        //                 `Server responded with status: ${response.status}, body: ${errorData}`
        //             )
        //         } else {
        //             const data = await response.text()
        //             console.log(
        //                 `${this.filesToUpload.length} files were transferred successfully!, body: ${data}`
        //             )
        //         }
        //     } catch (error) {
        //         console.error("Failed to upload files:", error.message)
        //     } finally {
        //         this.filesToUpload = []
        //         this.isUploading = false
        //     }
        // },
    }))

    Alpine.data("downloadState", () => ({
        filesToDownload: [],
        ws: null,

        listen() {
            console.log("start websocket on host:", window.location.host)
            this.ws = new WebSocket(`ws://${window.location.host}/ws/files`)
            this.ws.onmessage = (event) => {
                console.log("onmessage:", event)
                this.filesToDownload = JSON.parse(event.data) ?? []
                console.log("filesToDownload:", this.filesToDownload)
            }
            this.ws.onclose = () => {
                console.log("onclose")
            }
        },

        async handleDownload(file) {
            file.downloading = true
            console.log("handleDownload", file.name, file.downloadLink)
            //
            // const response = await fetch(file.downloadLink);
            // const reader = response.body.getReader();
            // const contentLength = +response.headers.get('Content-Length');
            //
            // let receivedLength = 0; // received that many bytes at the moment
            // let chunks = []; // array of received binary chunks (comprises the body)
            // while (true) {
            //     const {done, value} = await reader.read();
            //
            //     if (done) {
            //         break;
            //     }
            //
            //     chunks.push(value);
            //     receivedLength += value.length;
            //
            //     console.log(`Received ${receivedLength} of ${contentLength}`)
            // }
            //
            // let chunksAll = new Uint8Array(receivedLength);
            // let position = 0;
            // for (let chunk of chunks) {
            //     chunksAll.set(chunk, position);
            //     position += chunk.length;
            // }
            // let result = new TextDecoder("utf-8").decode(chunksAll);

            // console.log("result:",result)

            // Convert to a blob and create a download link
            // const blob = new Blob([chunksAll]);
            // const url = URL.createObjectURL(blob);
            // const a = document.createElement('a');
            // a.href = url;
            // a.download = "largefile.bin"; // specify the download file name
            // document.body.appendChild(a);
            // a.click();
            // a.remove();
            //
            // URL.revokeObjectURL(url);
            // console.log("Download completed.");

            // fetch(file.downloadLink)
            //     .then(response => {
            //         if (!response.ok) throw new Error('Network response was not ok');
            //         const contentLength = response.headers.get('content-length');
            //         if (!contentLength) throw new Error('Content-Length header is missing');
            //
            //         const total = parseInt(contentLength, 10);
            //         let loaded = 0;
            //
            //         return new Response(new ReadableStream({
            //             start(controller) {
            //                 const reader = response.body.getReader();
            //                 read();
            //
            //                 function read() {
            //                     reader.read().then(({done, value}) => {
            //                         if (done) {
            //                             controller.close();
            //                             return;
            //                         }
            //                         loaded += value.length;
            //                         file.progress = Math.floor((loaded / total) * 100);
            //                         controller.enqueue(value);
            //                         self.files = [...self.files]; // Trigger Alpine.js reactivity
            //                         read();
            //                     }).catch(error => {
            //                         console.error('Download failed:', error);
            //                         file.downloading = false;
            //                     });
            //                 }
            //             }
            //         }));
            //     })
            //     .then(response => response.blob())
            //     .then(blob => {
            //         file.downloading = false;
            //         const url = window.URL.createObjectURL(blob);
            //         const a = document.createElement('a');
            //         a.href = url;
            //         a.download = file.name;
            //         document.body.appendChild(a);
            //         a.click();
            //         a.remove();
            //         window.URL.revokeObjectURL(url);
            //     })
            //     .catch(error => {
            //         console.error('Error:', error);
            //         file.downloading = false;
            //         file.progress = 0;
            //     });
        },
    }))
})
