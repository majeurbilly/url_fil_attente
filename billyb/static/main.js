// main.js

const affichage = document.getElementById("affichage");
let counter = 1;
// async = manière d'executer du code pendant qu'une operation en attente
async function fetchItems() {
    try {

        // requêtes réseau - demander des données à un serveur distant.
        const response = await fetch('http://localhost:8080/api/v1/items');
        console.log(response)

        // donnée brute json - id/vlue/lenght
        const data = await response.json();
        console.log('Données reçues :', data);

        // Initialisation des donnée recu, si vide tableau vide
        const items = Array.isArray(data) ? data : data.items || [];
        console.log(items.length)

        // afficher la liste des nom recu dans le backend
        affichage.innerHTML = `
            <h3>Liste des Items :</h3>
            <lo>
                ${items.map(item => `<li>${counter++} - Nom: ${item.value}</li>`).join('')}
            </lo>
        `;


    } catch (error) {
        console.error('Erreur dans la récupération des items:', error);
        affichage.innerHTML = `<p style="color: red;">Erreur : ${error.message}</p>`;
    }
}

// Appel de la fonction - s'execute a la volé
fetchItems();
