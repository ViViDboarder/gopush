alias pb-done='gopush "Done"'
alias pb-success='gopush "Success"'
alias pb-failure='gopush "Failure"'
function pb-notify() { $* && pb-success || pb-failure; }
