bq load \
--autodetect \
--replace \
--source_format=NEWLINE_DELIMITED_JSON \
boxscores.boxscoresNBA \
./json/testout.json
