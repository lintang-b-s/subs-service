package com.lintang.netflik.movieservice.broker.message;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Data
@Builder
public class UpdateActorMessage {
    private int id;
    private String name;
}
