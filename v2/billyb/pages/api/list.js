// pages/api/list.js

let list = []; // Liste initialisée vide au démarrage

export default function handler(req, res) {
    switch (req.method) {
        case 'GET':
            // Renvoyer la liste complète
            res.status(200).json({ success: true, data: list });
            break;

        case 'POST':
            // Ajouter un nouvel item
            const newItem = req.body.item;
            if (!newItem) {
                res.status(400).json({ success: false, message: 'Item is required' });
                return;
            }
            list.push(newItem);
            res.status(201).json({ success: true, data: list });
            break;

        case 'DELETE':
            // Supprimer un item spécifique
            const itemToDelete = req.body.item;
            list = list.filter(item => item !== itemToDelete);
            res.status(200).json({ success: true, data: list });
            break;

        default:
            res.setHeader('Allow', ['GET', 'POST', 'DELETE']);
            res.status(405).end(`Method ${req.method} Not Allowed`);
    }
}
