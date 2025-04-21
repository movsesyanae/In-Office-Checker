# In-Office-Checker
Golang Program to check whether you are in office. It does so by checking the connected
wifi, and comparing it with the configured office wifi name.

# Setup
## Building
To build the application use the following commands:
```
cd src
go build
```
## Adding to zshrc or bashrc or ...
```
export PATH="<path to in_office/src>:$PATH"
```
## Non-Default File Storage Path
The application saves data in `~/.my_tools/in_office/` by default. To change this,
set the command line varible `IN_OFFICE_PREFIX` to some other directory.
```
export IN_OFFICE_PREFIX=~/Documents/in_office/
```
## Application Setup
Run the following command to setup all necessary files:
```
in_office init clean
```
To reset the metadata, without formatting the saved days in office csv file,
ommit the *clean* argument:
```
in_office init
```
## Office Wifi Name Setup
To set the wifi name the application should be checking against,
use the following command:
```
in_office set office_wifi_name <some wifi name>
```
Note that all the words after `office_wifi_name` will be treated as part of
the wifi name.

## Automatic Running Setup
There is a sample plist file is provided which will run `in_office check`
every day at 10:32AM and 11:32AM. This will automatically check and update
the days in office saved data as necessary.

### Editing plist File
Line 10 should be updated to the absolute path of the application binary.
Lines 34 and 37 should be edited to update the log location. Enter the
absolute path.
### Copying and Creating the daemon
Copy the plist file into the applicable location:
```
cp sample_files/in_office.plist ~/Library/LaunchAgents/
```
Tell macOS about your Mac plist launchd file
```
cd ~/Library/LaunchAgents/
launchctl load in_office.plist
```
