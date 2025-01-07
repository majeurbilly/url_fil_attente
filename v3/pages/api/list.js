// Liste initialisée vide au démarrage
let list = [];

// Gestionnaire d'API pour Next.js
export default function handler(req, res) {
    // Ajouter les en-têtes CORS
    res.setHeader('Access-Control-Allow-Origin', '*'); // Remplacez '*' par votre domaine front-end pour plus de sécurité
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, DELETE, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

    // Vérifie la méthode HTTP
    switch (req.method) {
        case 'GET':
            // Renvoyer la liste complète
            res.status(200).json({ success: true, data: list });
            break;

        case 'POST':
            // Ajouter un nouvel item à la liste
            const newItem = req.body.item;
            if (!newItem) {
                res.status(400).json({ success: false, message: 'Item is required' });
                return;
            }
            list.push(newItem);
            res.status(201).json({ success: true, data: list });
            break;

        case 'DELETE':
            // Supprimer un item spécifique de la liste
            const itemToDelete = req.body.item;
            list = list.filter(item => item !== itemToDelete);
            res.status(200).json({ success: true, data: list });
            break;

        case 'OPTIONS':
            // Répondre aux requêtes OPTIONS
            res.status(200).end();
            break;

        default:
            // Méthode non supportée
            res.setHeader('Allow', ['GET', 'POST', 'DELETE', 'OPTIONS']);
            res.status(405).end(`Method ${req.method} Not Allowed`);
    }
}
