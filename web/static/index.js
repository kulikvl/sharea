"use strict"

/*================= Utils =================*/

function formatBytes(bytes) {
    if (bytes === 0) return "0 Bytes"
    const sizes = ["Bytes", "KiB", "MiB", "GiB", "TiB"]
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return parseFloat((bytes / Math.pow(1024, i)).toFixed(1)) + " " + sizes[i]
}

function formatDate(dateTimeStr) {
    const dateTime = new Date(dateTimeStr)
    const now = new Date()

    const hoursDiff = (now - dateTime) / (1000 * 60 * 60)

    if (hoursDiff < 24) {
        return "Today"
    } else if (hoursDiff < 48) {
        return "Yesterday"
    } else {
        return dateTime.toLocaleDateString("en-US", {
            month: "short",
            day: "2-digit",
            year: "numeric",
        })
    }
}

/*================= Web socket init =================*/

const ws = new WebSocket(`ws://${window.location.host}/ws/files`)
ws.onmessage = (event) => {
    // Broadcast the message to all Alpine components that are listening
    document.dispatchEvent(
        new CustomEvent("ws-message", {
            detail: JSON.parse(event.data),
        })
    )
}
ws.onclose = () => {
    console.log("WebSocket connection closed")
}

/*================= Alpine.js init =================*/

document.addEventListener("alpine:init", () => {
    Alpine.data("serverState", () => ({
        storageOccupiedBytes: 0,

        init() {
            document.addEventListener("ws-message", (event) => {
                const files = event.detail ?? []
                this.storageOccupiedBytes = files.reduce(
                    (acc, file) => acc + file.size,
                    0
                )
            })
        },
    }))

    Alpine.data("uploadState", () => ({
        filesToUpload: [],

        // Simple FSM I guess?
        uploadStatus: "ready", // "ready" | "uploading" | "uploaded" | "error"
        updateStatus() {
            switch (this.uploadStatus) {
                case "ready":
                    this.uploadStatus = "uploading"
                    break
                case "uploading":
                    if (this.uploadStats.error) {
                        this.filesToUpload = []
                        this.uploadStatus = "error"
                        break
                    }

                    this.uploadStats.uploadedCount++
                    if (
                        this.uploadStats.uploadedCount ===
                        this.uploadStats.totalCount
                    ) {
                        this.filesToUpload = []
                        this.uploadStats.totalTime = (
                            (new Date().getTime() -
                                this.uploadStats.startTime) /
                            1000
                        ).toFixed(2)
                        this.uploadStatus = "uploaded"
                    }
                    break
                case "error":
                case "uploaded":
                    this.uploadStatus = "ready"
                    break
            }
        },

        uploadStats: {
            error: null,
            startTime: null,
            totalTime: null,
            speed: 0,
            totalSize: 0,
            totalCount: 0,
            uploadedSize: 0,
            uploadedCount: 0,
        },

        updateStats(bytesUploaded) {
            this.uploadStats.uploadedSize += bytesUploaded

            const timeElapsed =
                (new Date().getTime() - this.uploadStats.startTime) / 1000
            this.uploadStats.speed = this.uploadStats.uploadedSize / timeElapsed

            const progressBar = document.getElementById("uploadProgressBar")
            progressBar.value =
                (this.uploadStats.uploadedSize / this.uploadStats.totalSize) *
                100
        },

        handleAdd(files) {
            const expectedLength = this.filesToUpload.length + files.length
            const allFiles = this.filesToUpload.concat(files)
            const actualLength = new Set(allFiles.map((file) => file.name)).size
            if (expectedLength !== actualLength) {
                console.error("Files with the same name are not allowed!")
                return
            }
            this.filesToUpload = allFiles
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
            let xhr = new XMLHttpRequest()

            let prevLoaded = 0
            xhr.upload.onprogress = (event) => {
                // console.log(
                //     `upload.onprogress: loaded: ${event.loaded}, total: ${event.total}, so progress is ${event.loaded / event.total}`
                //
                this.updateStats(event.loaded - prevLoaded)
                prevLoaded = event.loaded
            }

            xhr.onloadend = () => {
                if (xhr.status !== 202) {
                    console.error(
                        `Error while uploading file ${file.name}:`,
                        xhr.response
                    )
                    this.uploadStats.error = xhr.response
                }

                this.updateStatus()
            }

            // xhr.timeout = 10000 // timeout 10 seconds, wait for server response
            xhr.open("POST", `/api/upload/${file.name}`)
            xhr.send(file)
        },

        handleBatchUpload() {
            if (this.uploadStatus !== "ready") return

            this.uploadStats = {
                startTime: new Date().getTime(),
                speed: 0,
                totalSize: this.filesToUpload.reduce(
                    (acc, file) => acc + file.size,
                    0
                ),
                totalCount: this.filesToUpload.length,
                uploadedSize: 0,
                uploadedCount: 0,
            }

            this.updateStatus()

            this.filesToUpload.forEach((file) => this.handleUploadFile(file))
        },

        handleUploadMore() {
            if (
                this.uploadStatus !== "uploaded" &&
                this.uploadStatus !== "error"
            )
                return

            this.updateStatus()
        },
    }))

    Alpine.data("downloadState", () => ({
        filesToDownload: [],

        init() {
            document.addEventListener("ws-message", (event) => {
                this.filesToDownload = event.detail ?? []
            })
        },
    }))
})
