package com.lintang.netflik.movieservice.api.response;

import lombok.*;

@Setter
@Getter
@NoArgsConstructor
@Builder
@AllArgsConstructor
@ToString
public class Category {

    private int id;
    private String name;
}
