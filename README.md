# Railopedia

## Structure

| Codebase              |      Description          |
| :-------------------- | :-----------------------: |
| [web](packages/web)   |      Astro frontend       |
| [functions](packages/functions) |     Golang API  |

## Branches
- main -> staging branch for all pr
- prod -> what's running in prod

## Contributions
Railopedia is open to contributions, however I'd recommend you open an issue first so that we can discuss what you're working on and how it reflects the project.

## Why
As a university student living 4 hours drive from home, I often have to get the train when I go back home. The UK's train ticket market is dominated by Trainline,
however I find this to rarely be the cheapest option and end up spending ages searching for cheaper prices on alternative websites. There exists skyscanner and google flights
to simultaneously search all flight options, so why not one for trains ? With this project I aim to help not only students like me save money travelling back home but
the rest of the british population, as well as raise awareness of other ticketing providers. 

## About
Railopedia currently has support to search the following websites
- [Trainline](https://www.thetrainline.com)
- [Trainpal](https://www.mytrainpal.com)
- [Raileasy](https://new.raileasy.co.uk)
- [Traintickets](https://www.traintickets.com/?/)

I'm planning on adding additional sites later on during development.
The data is obtained via scrapes where I'm using [Colly](https://github.com/gocolly/colly) or [Rod](https://github.com/go-rod/rod)

## License
This work is licensed under a
[Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License][cc-by-nc-sa].

[cc-by-nc-sa]: http://creativecommons.org/licenses/by-nc-sa/4.0/