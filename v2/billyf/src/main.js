const affichage = document.getElementById("affichage");
const inputUser = document.getElementById("inputUser");
const submitButton = document.querySelector("button");
const url = "http://localhost:3000/api/list";

// Bouton de soumission
submitButton.onclick = handleSubmit;

// Charger les items à l'ouverture
fetchItems();

// Gestion de la soumission
function handleSubmit() {
    postName();
    fetchItems();
}

// Gestion de suppression
function handleDelete(itemId) {
    deleteItem(itemId);
}

// Récupérer et afficher les items
async function fetchItems() {
    try {
        const response = await fetch(url);
        const data = await response.json();
        const items = Array.isArray(data.data) ? data.data : [];
        affichage.innerHTML = `
            <h3>Noms au tableau :</h3>
            <ol>
                ${items.map((item, index) => `
                    <li>${item ?? 'Inconnu'}
                        <button onclick="handleDelete('${item}')" id="${index}" type="button" class="btn btn-danger btn-sm">X</button>
                    </li>
                `).join('')}
            </ol>
        `;
    } catch (error) {
        console.error('Erreur dans la récupération des items:', error);
        affichage.innerHTML = `<p style="color: red;">Erreur : ${error.message}</p>`;
    }
}

// Envoyer un nom
async function postName() {
    const name = inputUser.value.trim();
    if (!name) {
        alert('Veuillez entrer un nom.');
        return;
    }
    try {
        const data = { item: name };
        await fetch(url, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        inputUser.value = '';
    } catch (error) {
        console.error('Erreur lors de l\'envoi du nom:', error);
        alert(`Erreur : ${error.message}`);
    }
    fetchItems();
}

// Supprimer un élément
async function deleteItem(itemId) {
    try {
        await fetch(url, {
            method: 'DELETE',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ item: itemId })
        });
    } catch (error) {
        console.error('Erreur lors de la suppression:', error);
    }
    fetchItems();
}
