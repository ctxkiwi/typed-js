<?php

$outPath = __DIR__ . '/1.tjs';
$out = "";

for ($i = 0; $i < 1000; $i++) {
    $out .= "
struct Person$i {
    fn: string
    ln: string = \"a default $i\"
}

Person$i x$i = {
    fn: \"something$i\"
}
    ";
}

file_put_contents($outPath, $out);
