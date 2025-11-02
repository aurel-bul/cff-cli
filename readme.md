# CFF-CLI

Un outil CLI écrit en GO qui permet de consulter l'horaire CFF.

<p align="center">
  <img width="300" src="./img/showoff.png">
</p>

## Installation

1. Télécharger l'exécutable
    Se rendre sur la [page des Releases](https://gitlab.forge.hefr.ch/aurelien.bulliard/cff-cli/-/releases) pour télécharger la bonne archive selon votre OS/Architecture
2. Décompresser l'archive

## Utilisation

Pour l'instant, cff-cli implémente deux fonctionnalités:

### Afficher les prochains départs d'une gare

<details>
<summary>Linux/MacOS</summary>

```sh
./cff <gare>
```

Exemple:

```sh
./cff Fribourg
```

</details>
<details>
<summary>Windows</summary>

```cmd
.\cff.exe <gare>
```

Exemple:

```cmd
.\cff.exe Fribourg
```

</details>

### Calculer un trajet

<details>
<summary>Linux/MacOS</summary>

```sh
./cff trip <de> <à>
```

Exemple:

```sh
./cff trip Romont Fribourg
```

</details>
<details>
<summary>Windows</summary>

```cmd
.\cff.exe trip <de> <à>
```

Exemple:

```cmd
.\cff.exe trip Romont Fribourg
```

</details>
Exemple de trajet:

<p align="center">
  <img width="300" src="./img/trip.png">
</p>

## Bugs connus

- Les retards ne sont ni affichés ni pris en compte
- Le nombre de départs d'une gare est parfois faux si plusieurs trains partent à la même heure

## Crédits

Ce projet utilise l'[API Transport d'OpenData Suisse](https://transport.opendata.ch/), qui utilise les données fournies par [search.ch](https://search.ch/timetable/api/help), dépendantes de l'[horaire CFF](https://cff.ch).

Projet développé par Aurélien Bulliard originellement dans le cadre du cours de Programmation élégante en Go de la [Haute-école d'ingénérie et d'architecture de Fribourg](https://heia-fr.ch).
