#version 450
uniform float x;
out vec4 frag_colour;
void main() {
    frag_colour = vec4(cos(x), 0.5, sin(x), 1.0);
}
