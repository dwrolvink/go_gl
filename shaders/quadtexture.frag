#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D texture1;
uniform float x;

void main() {
    vec4 tex = texture(texture1, TexCoord);
    FragColor = vec4(tex.r - abs(x*0.5), tex.g - abs(x*0.5), tex.b, 1.0);
}