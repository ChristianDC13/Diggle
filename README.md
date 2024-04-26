# Diggle 🔍
This is a small experiment I conducted to learn how search engines work essentially. After spending months studying common algorithms for information retrieval, this is what I've achieved.

The project is written in Go and consists of the following parts:
- **Crawler**: Responsible for fetching web pages from the internet and saving their information in the database.
- **Ranker**: Calculates the pagerank of each page after the information has been saved.
- **Monitor**: A small module that observes the crawler's progress, which is useful if the crawler runs on multiple instances to streamline the process.
- **Searcher**: Handles the search and sorting logic, as well as providing the API and the graphical interface for performing searches.

I have documented the entire process and the theory behind this implementation in a video available on YouTube.

[Watch the video here!](https://youtu.be/PZzTxYZihAk)

I greatly appreciate any contributions or suggestions to the project.

Author: [Christian de la Cruz](https://github.com/ChristianDC13)
