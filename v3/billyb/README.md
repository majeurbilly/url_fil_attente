# Billyb

Billyb(ackend) est un backend ultra-simple construit en Next.js.

## Description
Le service contient une API simple permettant de gérer une liste accessible au endpoint `/api/list`.

- Par **GET**, il retourne le contenu de la liste.
- Par **POST**, il ajoute un élément dans la liste.
- Par **DELETE**, il supprime un élément spécifique de la liste.

---

![img_1.png](img_1.png)

---

## Stack
Le projet utilise les technologies suivantes :

- **Next.js** (API routes)
- **Node.js**
- **TailwindCSS** (pour les styles si un frontend est ajouté)
- **PostCSS**

---

## Usage

### Dev-local (sans installation)
Assurez-vous que vous avez Node.js installé sur votre machine.

1. Installez les dépendances :
   ```bash
   npm install
   ```

2. Lancez le serveur de développement :
   ```bash
   npm run dev
   ```

3. Accédez à votre API en visitant `http://localhost:3000/api/list`.

---

### Production
Pour déployer l'application en production :

1. Construisez l'application :
   ```bash
   npm run build
   ```

2. Lancez le serveur en mode production :
   ```bash
   npm start
   ```

3. L'API sera accessible à l'adresse `http://<votre-domaine>:3000/api/list`.

---

## Endpoints de l'API

### Envoyer un item dans la liste
Ajoutez un nouvel élément dans la liste avec la commande suivante :
```bash
curl -X POST http://localhost:3000/api/list -H "Content-Type: application/json" -d '{"item": "<NOM DE MON ITEM>"}'
```

### Voir la liste des éléments
Consultez tous les éléments présents dans la liste :
```bash
curl http://localhost:3000/api/list | jq
```

### Supprimer un élément
Supprimez un élément spécifique de la liste :
```bash
curl -X DELETE http://localhost:3000/api/list -H "Content-Type: application/json" -d '{"item": "<NOM DE L'ITEM>"}'
```

---

## Extensions utiles
Ajoutez ici des informations sur les extensions ou outils supplémentaires nécessaires (comme TailwindCSS pour un frontend).

---

## Fonctionnalités futures
Quelques fonctionnalités à implémenter dans les prochaines versions :

1. Gestion des utilisateurs.
2. Frontend intégré.
3. Intégration des métriques et d'un healthcheck.
4. Tests unitaires complets.

---

## Licence
Ajoutez une description de la licence ici.

---

Pour toute question ou contribution, contactez https://github.com/majeurbilly

