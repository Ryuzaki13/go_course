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

    if (AddControl && typeof AddControl === "function") {
        AddControl();
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


//===============
//===============
//===============

function reg(login, pass, name, role) {
    Send("POST", "/reg", {
        Login: login,
        Password: pass,
        Name: name,
        Role: role
    });
}

function login(login, pass) {
    Send("POST", "/login", {Login: login, Password: pass});
}

function logout() {
    Send("POST", "/logout");
}