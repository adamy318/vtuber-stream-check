#!/bin/python3

import requests
import time
import json
import pickle


def getLiverData(url):
    response = requests.get(url)
    data = response.json()['pageProps']['liver']
    return data


def main():

    urlPrefix = 'https://www.nijisanji.jp/_next/data/AssLCig6Qu-ghaLHbaxgh/members/'

    allLivers = ["kuzuha"]

    response = requests.get('https://www.nijisanji.jp/_next/data/AssLCig6Qu-ghaLHbaxgh/members/kuzuha.json')
    initialLiver = response.json()['pageProps']['liver']

    liversData = {"kuzuha": 
                    {
                        "name": initialLiver['name'],
                        "slug": initialLiver['slug'],
                        "affiliation": initialLiver['affiliation'],
                        "english_name": initialLiver['english_name'],
                        "youtube_ch": initialLiver['social_links']['youtube'],
                        "twitter": initialLiver['social_links']['twitter']
                    }
                }

    livers = response.json()['pageProps']['livers']['contents']
    for liver in livers:

        if liver['slug'] == "kuzuha":
            continue
        time.sleep(2)

        newData = {}

        url = urlPrefix+liver['slug']+".json"
        liverData = getLiverData(url)

        if liver['slug'] not in liversData:

            jpName = liver['name']
            slug = liver['slug']
            enName = liverData['english_name']

            affiliation = liverData['affiliation']
            socials = liverData['social_links']
            youtube = ""
            twitter = ""
            if 'youtube' in socials:
                youtube = socials['youtube']
            if 'twitter' in socials:
                twitter = socials['twitter']

            newData = { 
                slug: 
                {
                    "name": jpName,
                    "slug": slug,
                    "english_name": enName,
                    "affiliation": affiliation,
                    "youtube_ch": youtube,
                    "twitter": twitter
                }
            }

            allLivers.append(slug)
            liversData[slug] = newData
    
    with open("livers.txt", "w") as f:
        for liver in allLivers:

            f.write(liver+":")
            json.dump(liversData[liver], f)
            if liver != allLivers[-1]:
                f.write('\n')



if __name__ == "__main__":
    main()