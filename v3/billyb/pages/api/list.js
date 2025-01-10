// Liste initialisée vide au démarrage
let list = [];

// Exportation de la fonction handler pour gérer les requêtes HTTP
export default function handler(req, res) {
    // Vérifie le type de méthode HTTP utilisée pour la requête
    switch (req.method) {
        case 'GET':
            // Si la méthode est GET, renvoyer la liste complète
            res.status(200).json({ success: true, data: list });
            break;

        case 'POST':
            // Si la méthode est POST, ajouter un nouvel item à la liste
            const newItem = req.body.item; // Récupérer le nouvel item envoyé dans la requête
            if (!newItem) { // Vérifie si un item est fourni dans la requête
                // Retourne une erreur 400 si l'item est manquant
                res.status(400).json({ success: false, message: 'Item is required' });
                return; // Stoppe l'exécution du code
            }
            // Ajoute l'item à la liste
            list.push(newItem);
            // Retourne la liste mise à jour avec un statut 201 (créé avec succès)
            res.status(201).json({ success: true, data: list });
            break;

        case 'DELETE':
            // Si la méthode est DELETE, supprimer un item spécifique de la liste
            const itemToDelete = req.body.item; // Récupérer l'item à supprimer depuis la requête
            // Filtrer la liste pour enlever l'item spécifié
            list = list.filter(item => item !== itemToDelete);
            // Retourne la liste mise à jour avec un statut 200 (succès)
            res.status(200).json({ success: true, data: list });
            break;

        case 'OPTIONS':
            res.status(200).json({ success: true, data: list });
            break;

        default:
            // Si une méthode HTTP non supportée est utilisée
            // Définit les méthodes autorisées dans l'en-tête de la réponse
            res.setHeader('Allow', ['GET', 'POST', 'DELETE', 'OPTIONS']);
            // Retourne une erreur 405 (méthode non autorisée)
            res.status(405).end(`Method ${req.method} Not Allowed`);
    }
}