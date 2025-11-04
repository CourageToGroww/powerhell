package modules

// GetLessonContent returns detailed content for a specific lesson
func GetLessonContent(moduleID, lessonID string) string {
	content := map[string]map[string]string{
		"basics": {
			"basics-1": `
## What is PowerShell?

PowerShell is a **cross-platform** task automation solution made up of:
- A command-line shell
- A scripting language
- A configuration management framework

### Why PowerShell?

1. **Object-Based**: Unlike traditional shells that work with text, PowerShell works with .NET objects
2. **Discoverable**: Built-in help system and predictable verb-noun syntax
3. **Powerful**: Access to .NET framework and Windows Management Instrumentation (WMI)
4. **Cross-Platform**: Runs on Windows, Linux, and macOS

### Your First Command

Try this simple command:
` + "```powershell\nGet-Process\n```" + `

This lists all running processes on your system!`,

			"basics-2": `
## Variables and Data Types

Variables in PowerShell are incredibly flexible containers for storing data.

### Basic Variable Assignment

` + "```powershell\n$name = \"Alice\"\n$age = 30\n$isActive = $true\n$items = @(1, 2, 3, 4, 5)\n```" + `

### Data Types

PowerShell supports many data types:
- **String**: Text data
- **Int/Double**: Numeric data
- **Boolean**: True/False values
- **Array**: Collections of items
- **HashTable**: Key-value pairs

### Type Casting

You can explicitly define types:
` + "```powershell\n[string]$text = \"42\"\n[int]$number = \"42\"\n[datetime]$date = \"2024-01-01\"\n```" + `

### Special Variables

PowerShell has several automatic variables:
- ` + "`$_`" + `: Current object in pipeline
- ` + "`$?`" + `: Success status of last command
- ` + "`$Error`" + `: Array of recent errors`,

			"basics-3": `
## Essential PowerShell Cmdlets

PowerShell cmdlets follow a **Verb-Noun** pattern, making them easy to understand and remember.

### File System Navigation

` + "```powershell\n# List items in current directory\nGet-ChildItem\n\n# Change directory\nSet-Location C:\\Users\n\n# Get current location\nGet-Location\n```" + `

### Working with Objects

` + "```powershell\n# Get properties of an object\nGet-Process | Get-Member\n\n# Select specific properties\nGet-Process | Select-Object Name, CPU\n\n# Filter objects\nGet-Process | Where-Object {$_.CPU -gt 10}\n```" + `

### Getting Help

` + "```powershell\n# Get help for a cmdlet\nGet-Help Get-Process\n\n# Get examples\nGet-Help Get-Process -Examples\n\n# Update help files\nUpdate-Help\n```" + ``,
		},
		"active-directory": {
			"ad-1": `
## Active Directory PowerShell Module

The Active Directory module for PowerShell is essential for managing AD environments.

### Installing the Module

` + "```powershell\n# On Windows Server\nInstall-WindowsFeature RSAT-AD-PowerShell\n\n# On Windows 10/11\nAdd-WindowsCapability -Online -Name Rsat.ActiveDirectory.DS-LDS.Tools\n```" + `

### Importing the Module

` + "```powershell\nImport-Module ActiveDirectory\n```" + `

### Basic AD Cmdlets

- **Get-ADUser**: Retrieve user information
- **New-ADUser**: Create new users
- **Set-ADUser**: Modify user properties
- **Remove-ADUser**: Delete users
- **Get-ADGroup**: Retrieve group information
- **Add-ADGroupMember**: Add members to groups`,

			"ad-2": `
## Managing AD Users with PowerShell

### Creating a New User

` + "```powershell\nNew-ADUser -Name \"John Smith\" `\n  -GivenName \"John\" `\n  -Surname \"Smith\" `\n  -SamAccountName \"jsmith\" `\n  -UserPrincipalName \"jsmith@contoso.com\" `\n  -Path \"OU=Users,DC=contoso,DC=com\" `\n  -AccountPassword (ConvertTo-SecureString \"P@ssw0rd123\" -AsPlainText -Force) `\n  -Enabled $true\n```" + `

### Bulk User Creation

` + "```powershell\n# Import users from CSV\n$users = Import-Csv \"C:\\users.csv\"\n\nforeach ($user in $users) {\n    New-ADUser -Name $user.FullName `\n      -GivenName $user.FirstName `\n      -Surname $user.LastName `\n      -SamAccountName $user.Username `\n      -Department $user.Department `\n      -Title $user.Title `\n      -Enabled $true\n}\n```" + `

### Modifying Users

` + "```powershell\n# Change user properties\nSet-ADUser -Identity \"jsmith\" -Title \"Senior Developer\" -Department \"IT\"\n\n# Reset password\nSet-ADAccountPassword -Identity \"jsmith\" -Reset -NewPassword (ConvertTo-SecureString \"NewP@ss123\" -AsPlainText -Force)\n```" + ``,
		},
	}

	if moduleContent, exists := content[moduleID]; exists {
		if lessonContent, exists := moduleContent[lessonID]; exists {
			return lessonContent
		}
	}
	
	return "Lesson content not available yet."
}

// GetExerciseForLesson returns an exercise for a specific lesson
func GetExerciseForLesson(moduleID, lessonID string) Exercise {
	exercises := map[string]map[string]Exercise{
		"basics": {
			"basics-1": {
				ID:           "ex-basics-1",
				Instructions: "Write a PowerShell command to get all processes that start with 'chrome' and display only their Name and CPU usage.",
				StarterCode:  "# Your code here\nGet-Process | ",
				Solution:     "Get-Process chrome* | Select-Object Name, CPU",
				Hints: []string{
					"Use Get-Process with a wildcard pattern",
					"Pipe the results to Select-Object",
					"The wildcard pattern should be 'chrome*'",
				},
			},
			"basics-2": {
				ID:           "ex-basics-2",
				Instructions: "Create a variable that stores an array of your favorite programming languages, then add 'PowerShell' to it.",
				StarterCode:  "# Create an array of languages\n$languages = @()\n\n# Add PowerShell to the array\n",
				Solution:     "$languages = @('Python', 'JavaScript', 'Go')\n$languages += 'PowerShell'",
				Hints: []string{
					"Use @() to create an array",
					"Use += to add an item to an array",
					"Don't forget the quotes around string values",
				},
			},
		},
	}

	if moduleExercises, exists := exercises[moduleID]; exists {
		if exercise, exists := moduleExercises[lessonID]; exists {
			return exercise
		}
	}

	return Exercise{
		ID:           "default",
		Instructions: "Practice what you've learned by experimenting with the commands shown in this lesson.",
		StarterCode:  "# Type your PowerShell commands here\n",
		Hints:        []string{"Review the lesson content", "Try the example commands"},
	}
}