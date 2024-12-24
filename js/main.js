const affichage = document.getElementById("affichage");

fetch("http://localhost:8080/api/v1/items").then(async function (response) {
    affichage.textContent = await response.json();
    return response.json();
}).then(function(data) {
    console.log(data);
}).catch(function(err) {
    console.log('Fetch Error :-S', err);
});
