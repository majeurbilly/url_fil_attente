// Code amélioré

const affichage = document.getElementById("affichage");
const inputUtilisateur = document.getElementById("inputUtilisateur");
const boutonSoumettre = document.querySelector("button");
const messageErreur = document.getElementById("messageErreur");
const url = "https://10.100.2.130:3000/web/";


addEventListener("DOMContentLoaded", VerificationConnextionBackend);


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
    messageErreur.textContent += " - ";
}

// Charger et afficher les éléments
async function ChargerElements() {
    try {
        const reponse = await fetch(url, {credentials: 'include'});
        const donnees = await reponse.json();
        const elements = Array.isArray(donnees.data) ? donnees.data : [];
        if (reponse.ok && elements.length > 0) {
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
        } else if (reponse.ok && elements.length === 0) {
            <h3>Noms au tableau :</h3>
            <h4>Tableau vide</h4>
        }

    } catch (erreur) {
        console.error("Erreur lors de la récupération des éléments :", erreur);
        affichage.innerHTML = `<p style="color: red;">Erreur : ${erreur.message}</p>`;
        DisplayError("Erreur lors de la récupération des éléments")
    }
}
