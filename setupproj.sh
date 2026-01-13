#!/usr/bin/env bash

git init
git add .
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:thepseudocodeus/signal-from-noise.git
git push -u origin main
