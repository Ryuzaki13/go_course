let input = document.querySelector('input[type="file"]');
let isPopupOpen = false;

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
        if (callback && typeof callback === "function") {
            callback(JSON.parse(this.response));
        }
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
    let divImage = document.createElement("div");
    let image = document.createElement("img");



    div2.textContent = product.name;
    div2.className = "caption";
    div3.textContent = product.price;
    div3.className = "price";
    div.dataset.id = product.id;
    div.className = "product";
    divImage.append(image);
    div.append(divImage, div2, div3);

    if (window["IsAdmin"]) {
        div.append(Div({
            className: "fa-solid fa-ellipsis-vertical product-menu",
            events: {onclick: openMenu}
        }));
    }

    return div;
}

function Div(props) {
    let div = document.createElement("div");

    if (!props || typeof props !== "object") {
        return div;
    }

    if ("events" in props) {
        for (let key in props.events) {
            div[key] = props.events[key];
        }
    }

    if ("className" in props) {
        div.className = props.className;
    }

    if ("style" in props) {
        div.style.cssText = props.style;
    }

    if ("id" in props) {
        div.id = props.id;
    }

    if ("dataset" in props) {
        for (let key in props.dataset) {
            div.dataset[key] = props.dataset[key];
        }
    }

    if ("children" in props) {
        div.append(...props.children);
    }

    return div;
}

function openMenu() {
    let product = this.closest(".product");

    removePopups();

    if (product) {
        let popup = product.querySelector(".popup");
        if (!popup) {
            product.append(Div({
                className: "popup",
                children: [
                    Div({
                        className: "fa-solid fa-pen"
                    }),
                    Div({
                        className: "fa-solid fa-trash-can",
                        dataset: {id: product.dataset.id},
                        events: {
                            onclick: removeProduct
                        }
                    })
                ]
            }));
            isPopupOpen = true;
        } else {
            popup.remove();
        }
    }
}

function removeProduct() {
    let product = this.closest(".product");
    Send("DELETE", "/api/product", {id: +this.dataset.id}, () => {
        if (product) {
            product.remove();
        }
    });
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
        timeout = setTimeout(() => {
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


document.body.addEventListener("click", event => {
    if (event.target.classList.contains("product-menu")) {
        return;
    }

    removePopups();
});

function removePopups() {
    if (isPopupOpen === true) {
        let popup = document.querySelectorAll(".product .popup");
        for (let i = 0; i < popup.length; i++) {
            popup[i].remove();
        }
        isPopupOpen = false;
    }
}


//===============
//===============
//===============

function reg() {
    Send("POST", "/reg", {Login: "admin", Password: "admin", Name: "admin", Role: "admin"});
}

function login() {
    Send("POST", "/login", {Login: "admin", Password: "admin"});
}

function logout() {
    Send("POST", "/logout");
}