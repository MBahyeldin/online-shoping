import { Hero } from '@/components/home/Hero'
import { FeaturedCakes } from '@/components/home/FeaturedCakes'
import { Testimonials } from '@/components/home/Testimonials'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Cake, Truck, Heart, Award } from 'lucide-react'

const features = [
  {
    icon: Cake,
    title: 'Custom Designs',
    description: 'Every cake is crafted to your exact specifications — flavors, sizes, and decorations.',
  },
  {
    icon: Heart,
    title: 'Made With Love',
    description: 'Our bakers pour their heart into every creation using premium, locally-sourced ingredients.',
  },
  {
    icon: Truck,
    title: 'Fresh Delivery',
    description: 'Delivered fresh on your chosen date and time, right to your doorstep.',
  },
  {
    icon: Award,
    title: 'Award-Winning',
    description: 'Recognized as the best cake shop in the region for 5 consecutive years.',
  },
]

export function HomePage() {
  return (
    <main>
      <Hero />

      {/* Features */}
      <section className="py-16 bg-white">
        <div className="container mx-auto max-w-6xl px-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((f) => (
              <div key={f.title} className="text-center space-y-3">
                <div className="inline-flex items-center justify-center h-14 w-14 rounded-2xl bg-orange-100 text-primary">
                  <f.icon className="h-7 w-7" />
                </div>
                <h3 className="font-bold text-gray-900">{f.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{f.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <FeaturedCakes />
      <Testimonials />

      {/* CTA banner */}
      <section className="py-20 bg-primary text-white">
        <div className="container mx-auto max-w-3xl px-4 text-center">
          <h2 className="text-3xl md:text-4xl font-extrabold mb-4">
            Ready to Order Your Dream Cake?
          </h2>
          <p className="text-primary-foreground/80 mb-8 text-lg max-w-xl mx-auto">
            Browse our full collection of handcrafted cakes or contact us for a
            fully custom creation. Your perfect cake is just a few clicks away.
          </p>
          <Button
            size="lg"
            variant="secondary"
            className="text-primary font-bold text-base"
            asChild
          >
            <Link to="/products">Shop All Cakes</Link>
          </Button>
        </div>
      </section>
    </main>
  )
}
