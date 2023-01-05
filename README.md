# "Building that would be ridiculous"
## No, it'd be _redoculous_

Toy static site generator/documentation generator project to satisfy my curiosity. If you use this in production you're either brave or mad, i salute you.

Fed up of most generators using markdown + front matter which breaks linting tools, broke out the configuration and the content
into separate files so semantically + technically correct markdown can be used exclusively.

Chuck your config files into `content/`, i've left the sample files i used during development in there.

Chuck your markdown files into `includes/`, again i've left the sample files i used during development.

Built HTML files will magically appear in `build/`

If you want to modify the templates that the build markdown is sandwiched between, those files can be found in `template/`