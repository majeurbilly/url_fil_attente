const affichage = document.getElementById("affichage");
const bouton = document.getElementById("bouton");
const saisie = document.getElementById("saisie");

// Appel API BillyB existante
bouton.addEventListener("click", () => {
    const nom = saisie.value;
    if (nom) {
        httpGet(`http://localhost:8080/api/items?nom=${encodeURIComponent(nom)}`);
    } else {
        affichage.innerHTML = "Veuillez entrer un nom.";
    }
});

function httpGet(theUrl) {
    fetch(theUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error("Erreur : " + response.statusText);
            }
            return response.text();
        })
        .then(data => {
            affichage.innerHTML = data;
        })
        .catch(error => {
            affichage.innerHTML = error;
        });
}