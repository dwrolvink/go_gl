#version 450
uniform float x;
uniform float y;
uniform float t;
out vec4 frag_colour;
void main() {
    frag_colour = vec4(0.7 - x, sin(0.2*t), 0.7 - y, 1.0);
}
