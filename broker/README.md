Test with: go test ./internal/math

**TODO**:
[ ] Parser should be smart enough to understand trading, after and pre market hours.
[ ] Volume or price can be missing for a certain time, if it is so, take next or previous minute data as a substitute
[ ] Should compare volume with average volume of each window
[ ] Use ratios instead of differences???
[ ] Find last non null / 0 volume record and reference (last point) in Yahoo get ticker
[ ] Write tests for existing functionality
