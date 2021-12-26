# TedTalks-Scraper
This program let you create a file .CSV with all information from [TedTalks](https://www.ted.com/talks), including:

- Title
- Description
- Views (Number of Views)
- Author
- Date
- Tags
- Link for the Page on Ted.com


## The project

This project is build in [GoLang](https://go.dev/) using for the first time the [Colly](https://github.com/gocolly/colly) library to scrap all Talks from the Ted website.
I hope to improve the quality and quantity of the generated file, to create a better Dataset to use and maybe one day, upload it to Kaggle.

## TODO
This isn't perfect so here a list of things to improve/add:

- [ ] Retrive the Thumb Link and Add into the Dataset
- [ ] Retrive the Duration of the video
- [ ] Add (Maybe in an other file) the Descriptions of Authors
- [ ] Publish on Kaggle

## License
MIT