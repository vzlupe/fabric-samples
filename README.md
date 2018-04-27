
## Senior Capstone Project - Victor Zuniga

Project directory: fabcar
Chaincode: chaincode/fabcar/go/fabcar.go

To build the project:
```
git clone -b master https://github.com/vzlupe/fabric-samples
cd fabric-samples
git checkout master
```
Install platform binaries:
```
curl -sSL https://goo.gl/6wtTN5 | bash -s 1.1.0
```
Export PATH environment variable:
```
export PATH=<path to download location>/bin:$PATH
```
Navigate to the project:
```
cd fabcar
```
Within the project you will find a reset.sh script. Run this to set the project environment. This should take a minute or two.
```
./reset.sh
```
Once this completes, the system is ready to use. The files which contain the temperature readings are located in the datasets directory. To run the system edit the invoke.js file and on line 16 enter the file which should be read. Then run it like so:
```
node invoke.js
```
This process can take some time. With the stemp1in100.txt file it will take about 6 minutes and with any of the larger files it will take about 1 hour. This time is based on the frequency of transaction requests which is determined by the number on line 175 of invoke.js. As is this number is 4000 indicating 4 seconds per transaction. The system may be able to handle lower times or require higher times depending on machine specs. A number of variations of invoke.js using different input files have been included in the testProgs directory for use if desired. They use a higher lenght of time between transactions to ensure that the system functions properly given the long runtime of the files. Simply move the file into the main fabcar directory and run the same as invoke.js. The running of the system uses the logic in chaincode file mentiioned at the top of this README. Once the transactions have completed, the blockchain can be queried using query.js like so:
```
node query.js
```
This will output a result to the file specified on line 117 of query.js. My output has been included in the testResults directory.
