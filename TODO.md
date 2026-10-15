# TODO

## Fixes

- [ ] error buttons keine glows
- [X] Scrollbare listen sind zu hoch
- [X] Water View macht nix
- [X] Water Icon zu Gießkanne
- [X] Manage zu My Plants umbennen
- [X] Menu schließen
- [X] Legal Seiten im Browser öffnen
- [X] Doppeltes Emojis auf der Startseite
- [X] Grüne Border um Plant Card entfernen
- [X] Translations updaten
- [X] Cache löschen nach Gießen
- [X] Nach Create 1. manage speichern wird 404
- [X] zuürck + edit buttons im plant view sind out of area
- [ ] Wenn man auf eine Pflanze klickt, dann kommt man zu <https://my.water-me.app/plant/6a6a6d0f39a107ddfcee3805>, was ein HTTP 404 wirft. Das ist auf CF Workers als static gehostet

## Features

- [ ] Notification service
- [X] water view for all plants that need watering today
- [X] plant care history view
- [X] Ansicht Pflanze öffnen mit bearbeiten Button dann
- [X] Neue Pflanze erstellen Button im Menü und unten in der Plant List
- [X] Water View mit Plant List
- [X] Dark Mode aktivierbar machen (folgt Geräteeinstellung prefers-color-scheme)
- [ ] Gucken ob Handy online ist und nur dann bearbeiten/gießen erlauben
- [X] Real content on website
- [X] Use Stats Endpoints in Website
- [X] Legal Pages
- [X] Legal Links funktionieren auf Webapp nicht
- [ ] Multi Replica Notication Worker
- [ ] Alle i18n files in einen ordner
- [ ] Frontend für Notification Configs
- [ ] Endpoints nur für Historie bearbeiten usw
- [ ] Fade von Scrollable ist nicht unten wenn Scrollbar nicht komplette Höhe ist
- [X] Aus Feuchtigkeitmessungen soll ein Gießplan berechnet werden (manuell/automatisch)
- [X] Tägliche Aufgabenliste (Pflanze gießen / Feuchtigkeit messen)
- [ ] LastWatered/Fertizlited checken

## Future

- [ ] kinda read only mode to share access with family members or neighbors when going on vacation
- [ ] plant species database integration to auto fill care info
- [ ] calendar view for plant care tasks
- [ ] Füge eine Art Flags hinzu, dass eine Pflanze Repotting zb braucht. Dann soll man damit Tasks planen können, dass man zb in einer Woche pflanze xy umtopfen möchte und soll dann an dem tag dran erinnert werden oder kann das verschieben usw
- [ ] Log out button
- [ ] Delete Account Button
- [ ] Delete Plants
- [ ] Feuchtigkeitsmessungen rückwirkend mit Zeitstempel eintragen (Backfill UI)
- [ ] Chart-Darstellung des Feuchtigkeitsverlaufs (aktuell nur Liste)
- [ ] Sensor-/Home-Assistant-Integration für automatische Messwerte

In the plant card shown in "my garden", only show essential information of a plant. Show the name and species, properly formatted next to each othe. Show last watered and next water date appealingly formatted. Other information doesnt need to be shown in the card for now.
