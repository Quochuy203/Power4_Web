## 🧩 Puissance 4 Web – Jeu interactif en Go
## 🎮 Présentation
Ce projet est une implémentation web du célèbre jeu Puissance 4, développé en Golang avec une interface en HTML/CSS. Il permet à deux joueurs de s’affronter ou de jouer contre un robot doté de plusieurs niveaux d’intelligence. Le jeu intègre une mécanique originale de gravité inversée toutes les 10 manches, ainsi que des scénarios de victoire, match nul et rejouabilité.

## 🚀 Installation & Lancement
Prérequis :
- Go 1.20 ou supérieur
- Navigateur web moderne
Lancer le projet :
- git clone https://github.com/nom-du-repo/puissance4-web.git
- cd puissance4-web
- go run main.go
Accéder ensuite à: http://localhost:8080

## 🕹️ Règles du jeu
- Saisir les noms des joueurs ou choisir de jouer contre un robot.
- Sélectionner la difficulté : easy, normal, hard.
- Chaque joueur joue à tour de rôle en plaçant un jeton dans une colonne.
- Tous les 10 tours, la gravité s’inverse (les jetons tombent vers le haut).
- Le jeu se termine par une victoire (4 jetons alignés), un match nul ou une possibilité de rejouer.

## 🤖 Intelligence Artificielle

- Facile : choix aléatoire.
- Normal : détecte les coups gagnants ou bloque l’adversaire.
- Difficile : utilise une évaluation stratégique du plateau pour maximiser ses chances (extensible en Minimax).

## 🧱 Architecture du projet

```

├── main.go                  # Serveur principal et définition des routes
├── web/
│   ├── game.go              # Logique du jeu : grille, gravité, victoire
│   ├── ai.go                # Intelligence artificielle : choix des coups
│   ├── handler.go           # Gestion des routes et des interactions
│   └── templates/           # Fichiers HTML pour le rendu
│       ├── index.html       # Page d’accueil
│       ├── welcome.html     # Page de bienvenue après saisie des noms
│       ├── game.html        # Interface principale du jeu
│       ├── draw.html        # Page affichée en cas de match nul
│       ├── win.html         # Page affichée en cas de victoire
├── static/
│   └── style.css            # Feuille de style pour l’interface
└── README.md                # Documentation du projet
```

## ✅ Fonctionnalités clés
- Interface HTML/CSS dynamique
- Plusieurs niveaux de difficulté
- Gravité inversée toutes les 10 manches
- Gestion du tour par tour
- Détection de victoire et match nul
- Mise en évidence des jetons gagnants
- Architecture backend claire et maintenable

## 📚 Notes techniques
- Le rendu HTML est entièrement piloté par le backend Go via des templates.
- Aucune logique JavaScript n’est utilisée pour garantir la clarté et la maintenabilité.
- Le projet peut être enrichi avec des animations CSS ou une sauvegarde des parties.

## 👨‍💻 Auteur
-  PHAM Huy 
