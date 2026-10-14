# Go Reloaded

Projet Go de transformation de texte inspiré des exercices de parsing et de formatage de phrases.

## Objectif

Ce programme lit un fichier texte, applique plusieurs transformations automatiques et écrit le résultat dans un fichier de sortie.

Les fonctionnalités incluent :

- transformation de casse sur des mots ciblés
- conversion hexadécimale / binaire vers décimal
- correction des espaces autour de la ponctuation
- ajout d'espaces autour des guillemets
- correction des articles `a` / `an`

## Fonctionnement

Le programme attend exactement 2 arguments :

```bash
go run . <input.txt> <output.txt>
```

Exemple :

```bash
go run . test/test.txt test/out.txt
```

## Transformations supportées

### 1) Mise en forme des mots
Les commandes sont indiquées sous la forme :

```text
(word) (up)
(word) (low)
(word) (cap)
```

Exemples :

- `this is a TEST (low)` → `this is a test`
- `hello world (up)` → `HELLO WORLD`
- `welcome to the city (cap)` → `Welcome To The City`

### 2) Conversion de nombres

- `1E (hex)` → `30`
- `1010 (bin)` → `10`

### 3) Ponctuation

- `there ,is a problem !with punctuation` → `there, is a problem! with punctuation`
- `this is a sentence . another sentence` → `this is a sentence. another sentence`

### 4) Articles

Le programme corrige automatiquement les articles en fonction du mot suivant :

- `a apple` → `an apple`
- `a banana` → `a banana`
- `an hour` → `an hour`

## Structure du projet

```text
Go-reloaded/
├── main.go
├── processe.go
├── go.mod
├── test/
│   ├── test.txt
│   └── out.txt
└── README.md
```

## Exemple de fichier d'entrée

```text
this is a TEST (low)
hello world (up)
welcome to the city (cap)
this is VERY COOL (low, 2)
this should be VERY NICE (up, 2)
```

## Résultat attendu

Le programme génère une version corrigée et plus lisible du texte, en appliquant les règles demandées automatiquement.

## Lancer le projet

```bash
go run . test/test.txt test/out.txt
```

Puis consultez le fichier de sortie généré dans `test/out.txt`.
