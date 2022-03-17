console.log(" %c I'm admin! ", "background-color:red;color:white;font-weight:bold");

let isPopupOpen = false;

function AddControl() {
    let productList = document.querySelectorAll(".product-list > .product");
    for (let i = 0; i < productList.length; i++) {
        productList[i].append(Div({
            className: "fa-solid fa-ellipsis-vertical product-menu",
            events: {onclick: openMenu}
        }));
    }
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

function removePopups() {
    if (isPopupOpen === true) {
        let popup = document.querySelectorAll(".product .popup");
        for (let i = 0; i < popup.length; i++) {
            popup[i].remove();
        }
        isPopupOpen = false;
    }
}

function addProduct() {

}

document.body.addEventListener("click", event => {
    if (event.target.classList.contains("product-menu")) {
        return;
    }

    removePopups();
});

function AddCreateProductButton() {
    let productList = document.querySelector("#ProductList");
    if (!productList) return;

    productList.prepend(Div({
        textContent: "+",
        events: {onclick: addProduct}
    }));
}

