# gopush
alias pb-done='gopush "Done"'
alias pb-success='gopush "Success"'
alias pb-failure='gopush "Failure"'
function pb-notify
  [ $status = 0 ] ;and pb-success ;or pb-failure
end
