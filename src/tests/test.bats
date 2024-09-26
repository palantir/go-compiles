@test "2: Wait 2 seconds" {
    start=$(date +%s)
    sleep 2
    end=$(date +%s)
    [[ $((end - start)) -ge 2 ]]
}
