import random
ifile = open("datasets/", "r") #enter file to read from
ofile = open("datasets/", "w") #enter file to be created
count = 0
for line in ifile:
 	a = random.randint(1,10001)
 	if (a%1000 == 0):
 		count += 1
  		b = line.strip() + ",accept,reject\n"
 	else:
  		b = line.strip() + ",accept,accept\n"
 	ofile.write(b)

print count