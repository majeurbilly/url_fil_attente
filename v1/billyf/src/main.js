const affichage = document.getElementById("affichage");
const inputUser = document.getElementById("inputUser");
const submitButton = document.querySelector("button");
const url = "http://localhost:8080/api/v1/items";
const deleteButton = document.getElementById("${item.id}");

///////main///////
// Attachement de la fonction nommée au bouton de soumission
submitButton.onclick = handleSubmit;



// Fonction nommée pour gérer la soumission
function handleSubmit() {
    postName();
    fetchItems();
}
function handleDelete(itemId) {
    deleteItem(itemId);
    // fetchItems();
}

// Charger les items à l'ouverture
fetchItems();
//////////////////

// Fonction GET pour récupérer et afficher les items
async function fetchItems() {
    try {
        // requêtes réseau - demander des données à un serveur distant.
        const response = await fetch(url);
        console.log('Réponse du serveur:', response);

        // donnée brute json - id/vlue/lenght
        const data = await response.json();
        console.log('Données reçues :', data);

        // Initialisation des donnée recu, si vide tableau vide
        const items = Array.isArray(data) ? data : data.items || [];
        console.log('Nombre d\'items reçus:', items.length);


        // Afficher la liste des noms reçu depuis le backend
        affichage.innerHTML = `
            <h3>Noms aux tableaux :</h3>
            <ol>
                <span>${items.map((item) => `<li>${item.value ?? 'Inconnu'}</span>
                <button onclick="handleDelete('${item.id}')" id="${item.id}" type="button" class="btn btn-danger btn-sm">X</button></li>`).join('')}
            </ol>
        `;
    } catch (error) {
        console.error('Erreur dans la récupération des items:', error);
        affichage.innerHTML = `<p style="color: red;">Erreur : ${error.message}</p>`;
    }
}

// Fonction POST pour envoyer un nom
async function postName() {
    const name = inputUser.value.trim(); // Récupère et nettoie la saisie utilisateur

    try {
        // clé/valeur - necessaire pour la requete post
        const data = {value: name};
        console.log('Données envoyées :', data);

        // requêtes réseau - dire quelles données j'envoie au serveur distant.
        const response = await fetch(url, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        console.log("Nom envoyé avec succès:", name);
        // Réinitialise le champ de saisie
        inputUser.value = '';

    } catch (error) {
        console.error('Erreur lors de l\'envoi du nom:', error);
        alert(`Erreur : ${error.message}`);
    }
    fetchItems();
}
async function deleteItem(itemId) {
    try {
        const response = await fetch(`http://localhost:8080/api/v1/item/${itemId}`, {
            method: 'DELETE',
            headers: {'Content-Type': 'application/json'},
        });
        console.log(`Élément avec l'ID ${itemId} supprimé avec succès.`);

    } catch (error) {
        console.error('Erreur lors de la suppression:', error);
    }
    fetchItems();
}
