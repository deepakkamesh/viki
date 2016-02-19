#!/bin/bash
echo "shipping viki to chip..."

ssh chip@chip sudo pkill vikid
scp vikid chip@chip:~/viki/
scp objects.conf chip@chip:~/viki/
scp ../resources/* chip@chip:~/viki/resources/


