# Airport-codes
Converting airport codes to client-friendly text

This program converts complex airport codes to client-friendly data, using airport-lookup.csv file as source.

I order for the program to work you have to use code in the input file and run the program like this:
- go run . ./input.txt ./output.txt ./airport-lookup.csv

I will give you one example to play with. Paste this sentence to input file:
- Your flight departs from #HAJ, and your destination is ##EDDW.
- Now run the program as described above.
- Check output file for desired output.

You can change # and ## codes as you will. It will always convert them into desired form as long as you take the codes that are listed in the airport-lookup.csv file.

- Try and change either both of them or one of the codes and see the result. (run the program again with new codes)

Facinating, isn't it?

There's more.

You can also convert hard-to-read date and times into easy-to-read format. To try out this use following text as input.

- D(2022-05-09T08:07Z)
- Run the program. It should have converted this time to 09 May 2022. 

- Try more with other times and dates as you will



