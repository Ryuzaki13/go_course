let input = document.querySelector('input[type="file"]');

if (input) {
    input.onchange = function (e) {
        Upload(this.files, (res) => {
            console.log(res);
        });
    }
}

function Upload(files, callback) {
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

function Send(method, uri, data, callback) {
    let xhr = new XMLHttpRequest();
    xhr.open(method, uri);

    xhr.onload = function (event) {
        callback(JSON.parse(this.response));
    }
    if (data) {
        console.log(data);
        xhr.setRequestHeader("Content-Type", "application/json; charset=utf-8");
        xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
        xhr.send(JSON.stringify(data));
    } else {
        xhr.send();
    }
}

Send("POST", "/api/product", {Search: ""}, response => {

    let productList = document.querySelector("#ProductList");
    if (!productList || !response) {
        return;
    }

    for (let product of response) {
        productList.append(createProduct(product));
    }
});

function createProduct(product) {
    let div = document.createElement("div");
    let div2 = document.createElement("div");
    let div3 = document.createElement("div");
    let del = document.createElement("div");
    del.className = "fa-solid fa-trash-can";

    del.onclick = removeProduct.bind(del, product.id);

    div2.textContent = product.name;
    div3.textContent = product.price;
    div.dataset.id = product.id;
    div.className = "product";
    div.append(div2, div3, del);

    return div;
}

function removeProduct(id) {
    Send("DELETE", "/api/product", {id}, ()=> {
        let product = this.closest(".product");
        if (product) {
            product.remove();
        }
    })
}

/**
 * @type {HTMLInputElement}
 */
let search = document.querySelector("#Search");
if (search) {
    let timeout = 0;
    search.oninput = function (e) {
        if (timeout !== 0) {
            clearTimeout(timeout);
            timeout = 0;
        }
        timeout = setTimeout(()=> {
            Send("POST", "/api/product", {Search: this.value}, response => {
                let productList = document.querySelector("#ProductList");
                if (!productList || !response) {
                    return;
                }

                productList.innerHTML = "";

                for (let product of response) {
                    productList.append(createProduct(product));
                }
            });
        }, 400);
    }
}