// Code amélioré

const affichage = document.getElementById("affichage");
const inputUtilisateur = document.getElementById("inputUtilisateur");
const boutonSoumettre = document.querySelector("button");
const messageErreur = document.getElementById("messageErreur");
const url = "http://10.100.2.130:3000/api/list";



addEventListener("DOMContentLoaded", VerificationConnextionBackend);
if (VerificationConnextionBackend) {
    ChargerElements();
    if (inputUtilisateur.value != null) {
        boutonSoumettre.onclick = GererSoumission;
    }
    else if (messageErreur.value == null) {
        DisplayError("inputUtilisateur est null")
    }

    document.addEventListener("DOMContentLoaded", ChargerElements);
    affichage.onclick = async function (event) {
        await GestionDelete(event);
    };
}


async function VerificationConnextionBackend() {
    try {
        const response = await fetch(url, {credentials: 'include'});
        if (response.ok) {
            console.log("Connexion au backend réussie.");
            return true;
        } else {
            console.error(`Problème de connexion : ${response.status} ${response.statusText}`);
            return false;
        }
    } catch (erreur) {
        console.error("Erreur lors de la connexion au backend :", erreur);
        DisplayError("Problème de connexion au backend", erreur);
        return false;
    }
}

function DisplayError(string) {
    messageErreur.textContent += string;
    messageErreur.classList.remove("d-none");
    messageErreur.textContent += " + ";
}

// Charger et afficher les éléments
async function ChargerElements() {
    try {
        const reponse = await fetch(url, {credentials: 'include'});
        const donnees = await reponse.json();
        const elements = Array.isArray(donnees.data) ? donnees.data : [];
        if (reponse.ok) {
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
        }
    } catch (erreur) {
        console.error("Erreur lors de la récupération des éléments :", erreur);
        affichage.innerHTML = `<p style="color: red;">Erreur : ${erreur.message}</p>`;
        DisplayError("Erreur lors de la récupération des éléments")
    }
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

// Envoyer un nom au backend
async function EnvoyerNom(nom) {
    try {
        const donnees = {item: nom};
        console.log(donnees);
        await fetch(url, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(donnees),
            credentials: 'include'
        });
        await ChargerElements();
    } catch (erreur) {
        throw new Error(erreur.message);
    }
}
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
// Supprimer un élément dans le backend
async function SupprimerElement(element) {
    try {
        const donnees = {item: element};
        await fetch(url, {
            method: "DELETE",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(donnees),
            credentials: 'include'
        });
        await ChargerElements();
    } catch (erreur) {
        throw new Error(erreur.message);
    }
}
