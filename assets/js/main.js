let input = document.querySelector('input[type="file"]');
if (input) {
    input.onchange = function (e) {
        Upload(this.files, (res) => {
            console.log(res);
        });
    }
}

function ClearElement(element) {
    if (!element || !(element instanceof HTMLElement)) return;

    while (element.children.length > 0) element.children[0].remove();
}

function Div(props) {
    let div = document.createElement("div");

    if (!props || typeof props !== "object") {
        return div;
    }

    if ("textContent" in props) {
        div.textContent = props.textContent;
    }

    if ("events" in props && typeof props.events === "object") {
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
                Product.Build(response);
            });
        }, 400);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    Send("POST", "/api/product", {Search: ""}, response => {

        Product.Build(response);

        // if (AddControl && typeof AddControl === "function") {
        //     AddControl();
        // }
    });
})

//===============
//===============
//===============

let Product = (function () {
    let object = {};

    let currentCount = 0;

    object.Build = function (products) {
        let productList = document.querySelector("#ProductList");
        if (!productList || !products) {
            return;
        }

        ClearElement(productList);

        for (let product of products) {
            productList.append(object.Create(product));
        }

        object.GetCart();
    }

    object.GetCart = function () {
        Send("GET", "/api/cart", null, r => {
            if (!r) return;

            let cartCompose = document.querySelector("#CartCompose");

            let count = 0;
            for (let cart of r) {
                count += cart.Count;
                let product = document.querySelector(`.product[data-id="${cart.ID}"] .counter > div:nth-child(2)`);
                if (product) {
                    product.textContent = cart.Count;
                }
                cartCompose.append(Div({
                    dataset: {id: cart.ID},
                    children: [
                        Div({
                            textContent: cart.Name
                        }),
                        Div({
                            textContent: cart.Price
                        }),
                        Div({
                            textContent: cart.Count
                        }),
                    ]
                }));
            }
            UpdateCartCount(count);
        })
    }

    object.Create = function (product) {
        let image = document.createElement("img");

        return Div({
            className: "product",
            dataset: {id: product.id},
            children: [
                Div({children: [image]}),
                Div({
                    className: "caption",
                    textContent: product.name,
                }),
                Div({
                    className: "price",
                    textContent: product.price,
                }),
                Div({
                    className: "counter",
                    children: [
                        Div({
                            className: "button",
                            textContent: "-",
                            events: {
                                onclick: object.Minus
                            }
                        }),
                        Div({textContent: "0"}),
                        Div({
                            className: "button",
                            textContent: "+",
                            events: {
                                onclick: object.Plus
                            }
                        }),
                        Div({
                            className: "button add",
                            textContent: "Добавить",
                            events: {
                                onclick: object.Add
                            }
                        })
                    ]
                })
            ]
        });
    }

    object.Minus = function (e) {
        let element = e.currentTarget;
        if (!element) return;

        if (element.nextElementSibling) {
            if (currentCount < 1) return;
            currentCount--;
            element.nextElementSibling.textContent = currentCount;
        }
    };
    object.Plus = function (e) {
        let element = e.currentTarget;
        if (!element) return;

        if (element.previousElementSibling) {
            currentCount++;
            element.previousElementSibling.textContent = currentCount;
        }
    };
    object.Add = function (e) {
        let element = e.currentTarget;
        if (!element) return;

        element = element.closest(".product");
        if (!element) return;

        let id = +element.dataset.id;

        Send("POST", "/api/cart/update", {
            ID: id,
            Count: currentCount,
        }, count => {
            UpdateCartCount(count);
        });
    }

    return object;
})();

let User = (function () {
    let object = {};

    object.Registration = (login, pass, name, role) => {
        Send("POST", "/user/reg", {
            Login: login,
            Password: pass,
            Name: name,
            Role: role
        });
    };

    object.Login = (login, pass) => {
        Send("POST", "/user/login", {Login: login, Password: pass});
    };

    object.Logout = () => {
        Send("POST", "/user/logout");
    };

    return object;
})();

function UpdateCartCount(count) {
    let cartCounter = document.querySelector(".cart-counter");
    if (!cartCounter) return;

    if (count > 0) {
        cartCounter.textContent = count;
    } else {
        cartCounter.textContent = "";
    }
}

