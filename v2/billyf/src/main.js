// Code amélioré

const affichage = document.getElementById("affichage");
const inputUtilisateur = document.getElementById("inputUtilisateur");
const boutonSoumettre = document.querySelector("button");
const messageErreur = document.getElementById("messageErreur");
const url = "https://IPPP/api/list";


addEventListener("DOMContentLoaded", VerificationConnextionBackend);


async function VerificationConnextionBackend() {
    try {
        const response = await fetch("https://IPPP/api/list"); //todo: Remplacez par une ip valide
        if (response.ok) {
            console.log("Connexion au backend réussie.");
            return true;
        } else {
            console.error(`Problème de connexion : ${response.status} ${response.statusText}`);
            DisplayError("Problème de connexion : " + response.statusText);
            return false;
        }
    } catch (erreur) {
        console.error("Erreur lors de la connexion au backend :", erreur);
        return false;
    }
}

// Initialiser l'application
if (VerificationConnextionBackend){
    boutonSoumettre.onclick = GererSoumission;
    document.addEventListener("DOMContentLoaded", ChargerElements);
}

function DisplayError(string) {
    messageErreur.textContent += string;
    messageErreur.classList.remove("d-none");
    messageErreur.textContent = " - ";
}

// Gérer la soumission du formulaire
async function GererSoumission() {
    const nom = inputUtilisateur.value.trim();
    console.log("nom :", nom);
    if (!nom) {
        DisplayError("Error veuillez entrer un nom")
        return;
    }
    try {
        await EnvoyerNom(nom);
        inputUtilisateur.value = '';
    } catch (erreur) {
        console.error("Error lorsque le nom a été envoyé :", erreur);
        alert(`Erreur : ${erreur.message}`);
        DisplayError("Error lorsque le nom a été envoyé" + "nom ")
    }
}

// Gérer les suppressions (via délégation d'événements)
affichage.onclick = async function (event) {
    await GestionDelete(event);
};

async function GestionDelete(event) {
    if (event.target.classList.contains("btn-danger")) {
        const element = event.target.dataset.element;
        try {
            await SupprimerElement(element);
            console.log("Element supprimer:", element);
        } catch (erreur) {
            console.error("Erreur lors de la suppression :", erreur);
            DisplayError("Error lors de la suppression");
        }
    }
}

// Charger et afficher les éléments
async function ChargerElements() {
    try {
        const reponse = await fetch(url);
        console.log(reponse);
        const donnees = await reponse.json();
        console.log(donnees);
        const elements = Array.isArray(donnees.data) ? donnees.data : [];
        console.log(elements);
        affichage.innerHTML = `
            <h3>Noms au tableau :</h3>
            <ol>
                ${elements.map((element) => `
                    <li>${element ?? "Inconnu"}
                        <button data-element="${element}" type="button" class="btn btn-danger btn-sm" aria-label="Supprimer ${element ?? 'Inconnu'}">X</button>
                    </li>
                `).join("")}
            </ol>
        `;
    } catch (erreur) {
        console.error("Erreur lors de la récupération des éléments :", erreur);
        affichage.innerHTML = `<p style="color: red;">Erreur : ${erreur.message}</p>`;
        DisplayError("Erreur lors de la récupération des éléments")
    }
}

// Envoyer un nom au backend
async function EnvoyerNom(nom) {
    try {
        const donnees = {item: nom};
        console.log(donnees);
        await fetch(url, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(donnees),
        });
        await ChargerElements();
    } catch (erreur) {
        throw new Error(erreur.message);
    }
}

// Supprimer un élément dans le backend
async function SupprimerElement(element) {
    try {
        const donnees = {item: element};
        await fetch(url, {
            method: "DELETE",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(donnees),
        });
        await ChargerElements();
    } catch (erreur) {
        throw new Error(erreur.message);
    }
}
