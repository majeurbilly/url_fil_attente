function ajouterNom() {
  // Récupérer le nom saisi
  const nom = document.getElementById('nom').value;

  // Créer un nouvel élément de liste (li)
  const li = document.createElement('li');
  li.textContent = nom;

  // Ajouter un bouton de suppression à côté du nom
  const boutonSupprimer = document.createElement('button');
  boutonSupprimer.textContent = 'X';
  boutonSupprimer.onclick = function() {
    li.remove(); // Supprimer l'élément de liste
  };
  li.appendChild(boutonSupprimer);

  // Ajouter l'élément de liste à la liste
  document.getElementById('liste').appendChild(li);

  // Effacer le champ de saisie
  document.getElementById('nom').value = '';
}