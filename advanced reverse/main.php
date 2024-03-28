<?php

// Задание:
// * Реализовать функцию, которая инвертирует порядок символов в каждом слове строки, 
//   оставляя при этом пунктуацию и пробелы на своих местах. Кроме того, функция сохраняет регистр символов, 
//   то есть в инвертированном слове сохраняется регистр исходного слова. 
// * Протестировать код, написав unit-тесты.

function advanced_reverse($sentence) {
    preg_match_all('/\w+\'*\w*|\p{P}+|\s+|`+/u', $sentence, $matches);

    $reversed_sentence = array();

    foreach ($matches[0] as $word) {
        if (preg_match('/^\p{P}+$/', $word) || $word === " ") { 
            $reversed_sentence[] = $word;
        } else {
            $reversed_word = '';
            $parts = preg_split("/('|`)/u", $word, -1, PREG_SPLIT_DELIM_CAPTURE);
            foreach ($parts as $index => $part) {
                if ($index % 2 === 0) {
                    preg_match_all('/./u', $part, $letters);
                    foreach ($letters[0] as $char) {
                        if (mb_strtolower($char, 'UTF-8') === $char) {
                            $reversed_word .= mb_strtolower(array_pop($letters[0]), 'UTF-8');
                        } else {
                            $reversed_word .= mb_strtoupper(array_pop($letters[0]), 'UTF-8');
                        }
                    }
                } else {
                    $reversed_word .= $part; 
                }
            }
            $reversed_sentence[] = $reversed_word;
        }
    }

    return implode('', $reversed_sentence);
}


function test_reverse() {
    $test_cases = array(
        array("Cat", "Tac"),
        array("Мышь", "Ьшым"),
        array("houSe", "esuOh"),
        array("домИК", "кимОД"),
        array("elEpHant", "tnAhPele"),
        array("cat,", "tac,"),
        array("Зима:", "Амиз:"),
        array("is 'cold' now", "si 'dloc' won"),
        array("это «Так» \"просто\"", "отэ «Кат» \"отсорп\""),
        array("third-part", "driht-trap"),
        array("can`t", "nac`t"),
        array("Die \"Küche\"", "Eid \"Ehcük\"")
    );

    foreach ($test_cases as $test) {
        $input = $test[0];
        $expected_output = $test[1];
        $result = advanced_reverse($input);
        if ($result === $expected_output) {
            echo "Test passed for input: $input\n";
        } else {
            echo "Test failed for input: $input\n";
            echo "Expected: $expected_output\n";
            echo "Got: $result\n";
        }
    }
}

test_reverse();
