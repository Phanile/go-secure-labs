const API_URL = "/api/products";


const productsContainer =
    document.getElementById("products");

const loadingElement =
    document.getElementById("products-loading");

const errorElement =
    document.getElementById("products-error");

const productsCountElement =
    document.getElementById("products-count");


async function fetchProducts() {
    const response = await fetch(API_URL);

    if (!response.ok) {
        throw new Error(
            `Failed to fetch products: ${response.status}`
        );
    }

    return response.json();
}


function createProductCard(product) {
    const card = document.createElement("article");

    card.className = "product-card";

    const imageWrapper =
        document.createElement("div");

    imageWrapper.className =
        "product-card__image-wrapper";


    const image =
        document.createElement("img");

    image.className =
        "product-card__image";

    image.src =
        product.image_url;

    image.alt =
        product.name;

    image.loading =
        "lazy";


    imageWrapper.appendChild(image);


    const content =
        document.createElement("div");

    content.className =
        "product-card__content";


    const name =
        document.createElement("h3");

    name.className =
        "product-card__name";

    name.textContent =
        product.name;


    const description =
        document.createElement("p");

    description.className =
        "product-card__description";

    description.textContent =
        product.description ?? "";


    const footer =
        document.createElement("div");

    footer.className =
        "product-card__footer";


    const price =
        document.createElement("span");

    price.className =
        "product-card__price";

    price.textContent =
        formatPrice(product.price);


    const button =
        document.createElement("button");

    button.className =
        "product-card__button";

    button.type =
        "button";

    button.textContent =
        "В корзину";


    button.addEventListener("click", () => {
        console.log(
            "Add product to cart:",
            product.id
        );
    });


    footer.appendChild(price);
    footer.appendChild(button);

    content.appendChild(name);
    content.appendChild(description);
    content.appendChild(footer);

    card.appendChild(imageWrapper);
    card.appendChild(content);

    return card;
}


function formatPrice(price) {
    return new Intl.NumberFormat(
        "ru-RU",
        {
            style: "currency",
            currency: "RUB",
            maximumFractionDigits: 2
        }
    ).format(price);
}


function renderProducts(products) {
    productsContainer.innerHTML = "";

    productsCountElement.textContent =
        `${products.length} товаров`;


    for (const product of products) {
        const card =
            createProductCard(product);

        productsContainer.appendChild(card);
    }
}


async function init() {
    try {
        const products =
            await fetchProducts();

        renderProducts(products);

        loadingElement.classList.add("hidden");

    } catch (error) {

        console.error(
            "Failed to load products:",
            error
        );

        loadingElement.classList.add("hidden");

        errorElement.classList.remove("hidden");
    }
}


init();