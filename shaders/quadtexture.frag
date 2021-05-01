#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D texture1;
uniform float x;
uniform float tex_x;
uniform float tex_y;
uniform float tex_divisions;
uniform float tex_fliph;

vec4 getSpriteTexture() {
    // flip
    vec2 n_texcoord;
    if (tex_fliph >= 1.0){
        n_texcoord = vec2(1.0 - TexCoord.s, TexCoord.t);
    }
    else {
        n_texcoord = TexCoord;
    }

    // Get the size of a tile on the entire texture
    float T = 1.0 / tex_divisions;

    // Get the tile coordinates
    float Tx = tex_x; // simple one
    float Ty = (tex_divisions - 1.0) - tex_y; // y:0 will be topleft, but is normally bottomleft

    return texture(
        texture1, 
        vec2(
            T * (n_texcoord[0] + Tx), 
            T * (n_texcoord[1] + Ty)
        )
    );
}

void main() {
    vec4 tex = getSpriteTexture();
    if (tex.a <= 0.1 ){
        discard;
    }
    FragColor = vec4(tex.r, tex.g, tex.b, tex.a);
}