let input = document.querySelector<HTMLInputElement>('input[type="file"]');

if (input) {
    input.onchange = function (this: HTMLInputElement, e: Event) {
        Upload(this.files, (res) => {
            console.log(res);
        });
    }
}

function Upload(files: FileList, callback: (r: any) => void) {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/upload");
    let data = new FormData();
    for (let file of files) {
        data.append("MyFiles", file, file.name);
    }
    xhr.onload = function (event) {
        callback(JSON.parse(this.response));
    }
    xhr.send(data);
}

