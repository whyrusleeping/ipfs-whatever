# ipfs-whatever
Some benchmarks for ipfs

## Usage

Run a version of the ipfs daemon you want to benchmark against, then:

```bash
$ ipfs-whatever > before.json
```

Then run your new patched version of ipfs and:

```bash
$ ipfs-whatever --before=before.json
checking patch operations per second...
checking 10MB file adds...
checking add-link ops per second...

Results
PatchOpsPerSec   1671.82  1763.93  5.51%
DirAddOpsPerSec  811.72   834.34   2.79%
Add10MBTime      78.75    81.96    4.07%
Add10MBStdev     16.11    13.70    -14.97%
```
