#!/bin/bash

HOOK_PATH=./.git/hooks/pre-commit

echoNl () {
    echo -en "\n# -------\n" >> $HOOK_PATH
}

echoToHook () {
    echo $1 >> $HOOK_PATH
    echoNl
}

catToHook () {
    cat $1 >> $HOOK_PATH
    echoNl
}


rm $HOOK_PATH
touch $HOOK_PATH

echoToHook "#!/bin/bash"
catToHook "./tools/precommit-hooks/lint.sh"
catToHook "./tools/precommit-hooks/prepare-commit-msg.sh"

chmod 755 ./.git/hooks/pre-commit